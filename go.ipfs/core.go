package ipfs

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"sync"

	config "github.com/ipfs/go-ipfs-config"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"

	"github.com/ipfs/go-ipfs/core"
	corehttp "github.com/ipfs/go-ipfs/core/corehttp"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	sockets "github.com/libp2p/go-socket-activation"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

/// ------ Setting up the IPFS Repo

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

/// ------ Spawning the node

func createConfig() (*config.Config, error) {
	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return nil, err
	}

	cfg.Experimental = config.Experiments{
		UrlstoreEnabled: true,
		QUIC:            true,
		PreferTLS:       true,
	}

	return cfg, nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, cfg *config.Config, repoPath string) (*core.IpfsNode, error) {
	// Create the repo with the config
	err := fsrepo.Init(repoPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init repo folder: %s", err)
	}

	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:    true,
		Permanent: true,
		Routing:   libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	node.IsDaemon = true

	// Attach the Core API to the constructed node
	return node, nil
}

// printSwarmAddrs prints the addresses of the host
func printSwarmAddrs(node *core.IpfsNode) {
	if !node.IsOnline {
		fmt.Println("Swarm not listening, running in offline mode.")
		return
	}

	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		fmt.Printf("failed to read listening addresses: %s", err)
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Strings(lisAddrs)
	for _, addr := range lisAddrs {
		fmt.Printf("Swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Strings(addrs)
	for _, addr := range addrs {
		fmt.Printf("Swarm announcing %s\n", addr)
	}

}

// serveHTTPGateway collects options, creates listener, prints status message and starts serving requests
func serveHTTPGateway(node *core.IpfsNode, cfg *config.Config) (<-chan error, error) {
	listeners, err := sockets.TakeListeners("io.ipfs.gateway")
	if err != nil {
		return nil, fmt.Errorf("serveHTTPGateway: socket activation failed: %s", err)
	}

	listenerAddrs := make(map[string]bool, len(listeners))
	for _, listener := range listeners {
		listenerAddrs[string(listener.Multiaddr().Bytes())] = true
	}

	gatewayAddrs := cfg.Addresses.Gateway
	for _, addr := range gatewayAddrs {
		gatewayMaddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPGateway: invalid gateway address: %q (err: %s)", addr, err)
		}

		if listenerAddrs[string(gatewayMaddr.Bytes())] {
			continue
		}

		gwLis, err := manet.Listen(gatewayMaddr)
		if err != nil {
			return nil, fmt.Errorf("serveHTTPGateway: manet.Listen(%s) failed: %s", gatewayMaddr, err)
		}
		listenerAddrs[string(gatewayMaddr.Bytes())] = true
		listeners = append(listeners, gwLis)
	}

	for _, listener := range listeners {
		fmt.Printf("Gateway (%s) server listening on %s\n", "readonly", listener.Multiaddr())
	}

	var opts = []corehttp.ServeOption{
		corehttp.MetricsCollectionOption("gateway"),
		corehttp.IPNSHostnameOption(),
		corehttp.GatewayOption(false, "/ipfs", "/ipns"),
		corehttp.VersionOption(),
		corehttp.CheckVersionOption(),
	}

	if cfg.Experimental.P2pHttpProxy {
		opts = append(opts, corehttp.ProxyOption())
	}

	if len(cfg.Gateway.RootRedirect) > 0 {
		opts = append(opts, corehttp.RedirectOption("", cfg.Gateway.RootRedirect))
	}

	errc := make(chan error)
	var wg sync.WaitGroup
	for _, lis := range listeners {
		wg.Add(1)
		go func(lis manet.Listener) {
			defer wg.Done()
			errc <- corehttp.Serve(node, manet.NetListener(lis), opts...)
		}(lis)
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	return errc, nil
}
