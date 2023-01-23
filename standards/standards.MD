# Standards

## What are Standards?
Standards are a generic concept that allows anyone to define interfaces and standards for whatever purpose they desire. All standards should be defined through minting a new BitBadge, but the purpose of that standard does not have to be related to BitBadges (it can be anything!).

## Genesis Standard (ID #0)
The genesis standard defines the base case with ID == 0. 

Details: All standards (including this one) must define some general things: 1) the details of the standard, and 2) which standards it extends (aka parent standards), and 3) author information. This genesis standard is the only standard that should not have any parents. 

The details of the standard can be defined as anything as long as it obeys the parent standard's minimum requirements. Additional details can be added as long as they don't interfere.

All standards must be instantiated by minting a new BitBadge where the implemented standard == 0 (genesis standard). For the genesis standard, we don't require any certain BitBadge permissions or metadata formats. Child standards can further require this. 

Standards are also recommended to be version controlled. 

Parent Standards: Since this is the genesis standard, we do not have any parent standards. This standard should be the only one without any parents.

Author: BitBadges Team

## URI + Bytes Standard (ID #1)
Details: This standard defines the general parent standard for all standards that defines a URI and arbitrary bytes field such as on BitBadges. 

Parent Standard: 0 (Genesis Standard)

Author: BitBadges Team

## BitBadge Standard (ID #2)
Details: This standard defines the parent standard for all BitBadges on the BitBadges blockchain. 

BitBadges are defined with the following fields stored on-chain: 
```
    uint64 id = 1; //Unique ID of the BitBadge
    UriObject uri = 2; //URI of the BitBadge and all badges
    bytes arbitraryBytes = 3; //Arbitrary bytes field (max 256 bytes)
    uint64 manager = 4; //Address number of the BitBadge manager
    uint64 permissions = 5; //Permissions (multiple bools packed into a uint64)
    repeated IdRange freezeRanges = 6; //Address numbers frozen or unfrozen
    uint64 nextSubassetId = 7; //Next badge ID to be used
    repeated BalanceObject BadgeSupplys = 8; //Subasset total supplies
    uint64 defaultSubassetSupply = 9; //Default subasset supply
    uint64 standard = 10; //Standard the badge implements
```

Child standards can further limit what these options can be on a BitBadge, such as what permissions are acceptable or who can be manager. Note there is no on-chain enforcement of these standards except for the fact that the standard ID# is stored with the badge on-chain. 

Any child standard must define how to verify any BitBadge that claims to implement it using the badge information. For example, for a standard that requires supply to be permanently locked, everyone must verify that the permissions are set to lock supply permanently, or else, it does not successfully implement this standard. Example child standards can be collectibles, collections, badges, soul bound tokens, etc. Child standards can even define different metadata interfaces for each type of BitBadge.

Parent Standard: 1 (URI + Bytes Standard)

Author: BitBadges Team

## Standard Badge w/ Public Metadata (ID #3)
TODO 
Categories, Tags, Other Metadata
    isWrapped : level it is wrapped
## Collectible Standard (ID #3)
TODO
## Proof of Attendance Badge Standard (ID #4)
TODO
## Collections Standard (ID #5)
TODO
Can have different collection types such as fully on-chain, trust but verify with URI JSON metadata, or fully private. Can also vary based on the implementation of the collection (merkle tree vs RSA accumulator).
## Verifiable Credential Badge Standard (ID #6)
                
