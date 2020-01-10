
package com.realtimeads.ipfs

import android.content.Context
import android.util.Log
import com.facebook.react.bridge.*
import java.io.*
import io.ipfs.Ipfs
import java.lang.*

operator fun File.get(path:String)=File(this,path)val Context.store get()=getExternalFilesDir(null)!!["ipfs"]

class RNIpfsModule(var reactContext:ReactApplicationContext):

ReactContextBaseJavaModule(reactContext) {
    var appContext: Context = reactContext.applicationContext

    var daemon = Thread(Runnable {
      try {
        Ipfs.startNode(appContext.store.toString())
      } catch (err: Exception) {
        Log.e("RNIpfsModule", "error starting node", err)
      }
    })

    init {
       RNIpfsLifecycle(this, reactContext)
       startNode()
    }

    override fun getName(): String {
        return "RNIpfs"
    }

    fun startNode() {
      if(!daemon.isAlive()) {
        Log.d("RNIpfsModule", "starting ipfs node")

        daemon.setDaemon(true)
        daemon.start()

        Log.d("RNIpfsModule", "started ipfs node")
      }
    }

    fun whenForeground() {
      startNode()
    }

    fun whenBackground() {
    }

    fun whenKilled() {
    }

    @ReactMethod
    public fun hasFile(multihash: String, promise: Promise) {
        Thread(Runnable {
            try {
                var multihash = Ipfs.hasFile(multihash)
                promise.resolve(multihash)
            } catch (err: Exception) {
                promise.resolve(false)
            }
        }).start()
    }

    @ReactMethod
    public fun addFromURL(url: String, promise: Promise){
        Thread(Runnable {
            try {
                var multihash = Ipfs.addFromURL(url, true)
                promise.resolve(multihash)
            } catch (err: Exception) {
                promise.reject(err)
            }
        }).start()
    }

}
