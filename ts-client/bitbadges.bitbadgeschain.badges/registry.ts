import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgNewCollection } from "./types/badges/tx";
import { MsgUpdateMetadata } from "./types/badges/tx";
import { MsgUpdateCollectionApprovedTransfers } from "./types/badges/tx";
import { MsgArchiveCollection } from "./types/badges/tx";
import { MsgMintAndDistributeBadges } from "./types/badges/tx";
import { MsgTransferBadges } from "./types/badges/tx";
import { MsgDeleteCollection } from "./types/badges/tx";
import { MsgUpdateUserApprovedTransfers } from "./types/badges/tx";
import { MsgUpdateCollection } from "./types/bitbadgeschain/badges/tx";
import { MsgUpdateUserPermissions } from "./types/badges/tx";
import { MsgUpdateCollectionPermissions } from "./types/badges/tx";
import { MsgUpdateManager } from "./types/badges/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/bitbadges.bitbadgeschain.badges.MsgNewCollection", MsgNewCollection],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateMetadata", MsgUpdateMetadata],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateCollectionApprovedTransfers", MsgUpdateCollectionApprovedTransfers],
    ["/bitbadges.bitbadgeschain.badges.MsgArchiveCollection", MsgArchiveCollection],
    ["/bitbadges.bitbadgeschain.badges.MsgMintAndDistributeBadges", MsgMintAndDistributeBadges],
    ["/bitbadges.bitbadgeschain.badges.MsgTransferBadges", MsgTransferBadges],
    ["/bitbadges.bitbadgeschain.badges.MsgDeleteCollection", MsgDeleteCollection],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateUserApprovedTransfers", MsgUpdateUserApprovedTransfers],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateCollection", MsgUpdateCollection],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateUserPermissions", MsgUpdateUserPermissions],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateCollectionPermissions", MsgUpdateCollectionPermissions],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateManager", MsgUpdateManager],
    
];

export { msgTypes }