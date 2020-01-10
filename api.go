package ipfs

import (
	"context"
	"fmt"
	"net/url"
	"time"

	filestore "github.com/ipfs/go-filestore"
	files "github.com/ipfs/go-ipfs-files"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	icore "github.com/ipfs/interface-go-ipfs-core"
	ioptions "github.com/ipfs/interface-go-ipfs-core/options"
	path "github.com/ipfs/interface-go-ipfs-core/path"
)

var ipfsNode *core.IpfsNode
var hasNode = false

// HasNode – checks if we already have a node on memory
func HasNode(v core.IpfsNode) bool {
	return hasNode
}

// GetNode – gets current node in memory
func GetNode() *core.IpfsNode {
	if hasNode {
		return ipfsNode
	}

	return nil
}

// StartNode – gets and starts a node in case it wasnt started yet
func StartNode(repoPath string) {
	if !hasNode {
		fmt.Println("-- Getting an IPFS node running -- ")

		ctx, _ := context.WithCancel(context.Background())

		fmt.Println("Spawning node on the provided repo path")

		if err := setupPlugins(repoPath); err != nil {
			panic(fmt.Errorf("error while setting up plugins: %s", err))
		}

		cfg, err := createConfig()
		if err != nil {
			panic(fmt.Errorf("failed to create ipfs config: %s", err))
		}

		node, err := createNode(ctx, cfg, repoPath)
		if err != nil {
			panic(fmt.Errorf("failed to spawn ipfs node: %s", err))
		}

		fmt.Println("IPFS node is running")

		printSwarmAddrs(node)

		// construct http gateway
		_, err = serveHTTPGateway(node, cfg)
		if err != nil {
			panic(fmt.Errorf("failed to spawn http gateway: %s", err))
		}

		ipfsNode = node
		hasNode = true
	}
}

// GetCoreAPI – returns the ipfs node wrapped with the core API interface
func GetCoreAPI() (icore.CoreAPI, error) {
	if !hasNode {
		return nil, fmt.Errorf("error getting CoreAPI: doesnt exist yet")
	}

	return coreapi.NewCoreAPI(&*ipfsNode)
}

// HasFile – checks if ipfs has the specific file
func HasFile(multihash string) (bool, error) {
	fmt.Println("-- HasFile: before has node -- ")

	if !hasNode {
		return false, fmt.Errorf("error checking multihash: node is off")
	}
	fmt.Println("-- HasFile: before GetCoreAPI -- ")

	api, _ := GetCoreAPI()
	fmt.Println("-- HasFile: before context.WithTimeout -- ")

	ctx, cancel := context.WithTimeout(ipfsNode.Context(), 5*time.Second)
	defer cancel()
	fmt.Println("-- HasFile: before Unixfs.Ls -- ")

	results, err := api.Unixfs().Ls(ctx, path.New(multihash))
	fmt.Println("-- HasFile: before Unixfs.Ls err -- ")

	if err != nil {
		return false, fmt.Errorf("error ls multihash: ", err)
	}
	fmt.Println("-- HasFile: before for results -- ")

	for link := range results {
		if link.Err != nil {
      return false, nil
			// panic(link.Err)
		}

		return true, nil
	}
	fmt.Println("-- HasFile: before return -- ")

	return false, nil
}

// AddFromURL – adds a file from an url
func AddFromURL(urlString string, pin bool) (string, error) {
	if !hasNode {
		return "", fmt.Errorf("error adding from url: node is off")
	}

	if !filestore.IsURL(urlString) {
		return "", fmt.Errorf("unsupported url syntax: %s", urlString)
	}

	api, _ := GetCoreAPI()

	ctx, cancel := context.WithCancel(ipfsNode.Context())
	defer cancel()

	opts := []ioptions.UnixfsAddOption{
		ioptions.Unixfs.Pin(pin),
		ioptions.Unixfs.CidVersion(0),
		ioptions.Unixfs.Nocopy(false),
	}

	url, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("error parsing url: ", err)
	}

	file := files.NewWebFile(url)
	pth, err := api.Unixfs().Add(ctx, file, opts...)
	if err != nil {
		return "", fmt.Errorf("error adding url: ", err)
	}

	return pth.Cid().String(), nil
}

// // PinNode pins an ipld node
// func PinNode(node *core.IpfsNode, nd ipld.Node, recursive bool) error {
// 	ctx, cancel := context.WithTimeout(node.Context(), PinTimeout)
// 	defer cancel()

// 	defer node.Blockstore.PinLock().Unlock()

// 	err := node.Pinning.Pin(ctx, nd, recursive)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "already pinned recursively") {
// 			return nil
// 		}
// 		return err
// 	}

// 	return node.Pinning.Flush()
// }

// // UnpinNode unpins an ipld node
// func UnpinNode(node *core.IpfsNode, nd ipld.Node, recursive bool) error {
// 	return UnpinCid(node, nd.Cid(), recursive)
// }
