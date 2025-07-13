# Extending the Approval (Advanced)

We try to offer as much as we can in the native interface; however, in edge cases, you may run into scenarios where you need more advanced functionality, for example, cross-approval functionality or access to other blockchain data / modules.

Before trying to extend the approval, please consider workarounds and design considerations. Many approvals that you think may need extension can be altered to fit into the native interfaces by thinking out of the box.

If you truly do need a more custom implementation, you will have to build your own smart contract with CosmWASM or an alternative self-implementation and call into the x/badges module. We refer you to the Tutorials section for how to do so.&#x20;

Please also let us know if whatever you are missing can be added natively. Our end goal is to have no one ever write a smart contract.

