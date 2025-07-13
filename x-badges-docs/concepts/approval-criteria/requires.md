# Requires

You also have the following options to further restrict who can transfer to who.

**requireToEqualsInitiatedBy, requireToDoesNotEqualsInitiatedBy**

**requireFromEqualsInitiatedBy, requireFromDoesNotEqualsInitiatedBy**

These are pretty self-explanatory. You can enforce that we additionally check if the to or from address equals or does not equal the initiator of the transfer.&#x20;

Note that this is bounded to the addresses in the respective lists for to, from, and initiatedBy (i.e. **toList, fromList, initiatedByList**).
