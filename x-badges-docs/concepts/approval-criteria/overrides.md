# Override User Level Approvals

As mentioned in the transferability page, the collection-wide approvals can override the user-level approvals. This is done via **overridesFromOutgoingApprovals** or **overridesToIncomingApprovals**.

If set to true, we will not check the user's incoming / outgoing approvals for the approved balances respectively. Essentially, it is **forcefully** transferred without needing user approvals. This can be leveraged to implement forcefully revoking a badge, freezing a badge, etc.

IMPORTANT: The Mint address has its own approvals store, but since it is not a real address, they are always empty. **Thus, it is important that when you define approvals from the Mint address, you always override the outgoing approvals of the Mint address.** Or else, the approval will not work.

* <pre class="language-json"><code class="lang-json"><strong>"fromListId": "Mint", //represents the list with the "Mint" addres
  </strong>...
  "approvalCriteria": {
    "overridesFromOutgoingApprovals": true
    ...
  }
  </code></pre>

<figure><img src="../../../../.gitbook/assets/image (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1) (1)  (18).png" alt=""><figcaption></figcaption></figure>
