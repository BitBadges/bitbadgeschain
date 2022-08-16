Parent standards - Each must include everything from the previous but can narrow it down (permanent or editable -> editable); can add new stuff 
All implement 0, parent standard 3
Generic Universal Standard - ID 0, Implements None, Parent None
    URI + Bytes Standard - ID 1, Implements 0, Parent 0
        BitBadge Standard (If standard == 0 on our chain, assume this) - ID 2, Implements 0, Parent 1 - On-chain verification -> Acceptable permissions list? Other badge checks such as manager == NULL address? May have some stuff for off-chain verification 
            Collectible - ID 3, Implements 0, Parent 2 
                BitBadges Metadata Collectible
            Badge - ID 4, Implements 0, Parent 2
                BitBadges Metadata Public Badge
                    Categories, Tags, Other Metadata
                     isWrapped : level it is wrapped
                Metadata can be encrypted, ZKPs, VCs
            Collections - multiple levels, private, merkle, accumulator, etc (actually enforce these on-chain) - ID 5, Implements 0, Parent 2
                BitBadges Metadata Collections
                    Collections can be (Address-Supply) to implement private amounts
        Social
        Education    
        Science