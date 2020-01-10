
# react-native-ipfs

## Getting started

`$ npm install react-native-ipfs --save`

### Mostly automatic installation

`$ react-native link react-native-ipfs`

### Manual installation


#### iOS

1. In XCode, in the project navigator, right click `Libraries` ➜ `Add Files to [your project's name]`
2. Go to `node_modules` ➜ `react-native-ipfs` and add `RNTesting.xcodeproj`
3. In XCode, in the project navigator, select your project. Add `libRNTesting.a` to your project's `Build Phases` ➜ `Link Binary With Libraries`
4. Run your project (`Cmd+R`)<

#### Android

1. Open up `android/app/src/main/java/[...]/MainActivity.java`
  - Add `import com.reactlibrary.RNIpfsPackage;` to the imports at the top of the file
  - Add `new RNIpfsPackage()` to the list returned by the `getPackages()` method
2. Append the following lines to `android/settings.gradle`:
  	```
  	include ':react-native-ipfs'
  	project(':react-native-ipfs').projectDir = new File(rootProject.projectDir, 	'../node_modules/react-native-ipfs/android')
  	```
3. Insert the following lines inside the dependencies block in `android/app/build.gradle`:
  	```
      compile project(':react-native-ipfs')
  	```

#### Windows
[Read it! :D](https://github.com/ReactWindows/react-native)

1. In Visual Studio add the `RNTesting.sln` in `node_modules/react-native-ipfs/windows/RNTesting.sln` folder to their solution, reference from their app.
2. Open up your `MainPage.cs` app
  - Add `using Testing.RNTesting;` to the usings at the top of the file
  - Add `new RNIpfsPackage()` to the `List<IReactPackage>` returned by the `Packages` method


## Usage
```javascript
import RNTesting from 'react-native-ipfs';

// TODO: What to do with the module?
RNTesting;
```
  