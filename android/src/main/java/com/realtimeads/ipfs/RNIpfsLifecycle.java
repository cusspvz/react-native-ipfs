package com.realtimeads.ipfs;

import android.content.Context;

import com.facebook.react.bridge.LifecycleEventListener;
import com.facebook.react.bridge.ReactApplicationContext;

public class RNIpfsLifecycle implements LifecycleEventListener {
  private RNIpfsModule mRNIpfs;
  private ReactApplicationContext mReactContext;
  private Context mApplicationContext;

  public RNIpfsLifecycle(RNIpfsModule rnIpfs, ReactApplicationContext reactContext) {
    this.mRNIpfs = rnIpfs;
    this.mReactContext = reactContext;
    this.mApplicationContext = reactContext.getApplicationContext();

    reactContext.addLifecycleEventListener(this);
  }

  public void onHostResume() {
    mRNIpfs.whenForeground();
  }

  public void onHostPause() {
    mRNIpfs.whenBackground();
  }

  public void onHostDestroy() {
    mRNIpfs.whenKilled();
  }
}
