using ReactNative.Bridge;
using System;
using System.Collections.Generic;
using Windows.ApplicationModel.Core;
using Windows.UI.Core;

namespace Testing.RNTesting
{
    /// <summary>
    /// A module that allows JS to share data.
    /// </summary>
    class RNTestingModule : NativeModuleBase
    {
        /// <summary>
        /// Instantiates the <see cref="RNTestingModule"/>.
        /// </summary>
        internal RNTestingModule()
        {

        }

        /// <summary>
        /// The name of the native module.
        /// </summary>
        public override string Name
        {
            get
            {
                return "RNTesting";
            }
        }
    }
}
