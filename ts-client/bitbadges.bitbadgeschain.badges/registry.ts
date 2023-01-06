import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgTransferBadge } from "./types/badges/tx";
import { MsgUpdateUris } from "./types/badges/tx";
import { MsgRegisterAddresses } from "./types/badges/tx";
import { MsgRequestTransferManager } from "./types/badges/tx";
import { MsgSetApproval } from "./types/badges/tx";
import { MsgUpdateBytes } from "./types/badges/tx";
import { MsgNewBadge } from "./types/badges/tx";
import { MsgNewSubBadge } from "./types/badges/tx";
import { MsgUpdatePermissions } from "./types/badges/tx";
import { MsgFreezeAddress } from "./types/badges/tx";
import { MsgTransferManager } from "./types/badges/tx";
import { MsgRevokeBadge } from "./types/badges/tx";
import { MsgRequestTransferBadge } from "./types/badges/tx";
import { MsgSelfDestructBadge } from "./types/badges/tx";
import { MsgHandlePendingTransfer } from "./types/badges/tx";
import { MsgPruneBalances } from "./types/badges/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/bitbadges.bitbadgeschain.badges.MsgTransferBadge", MsgTransferBadge],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateUris", MsgUpdateUris],
    ["/bitbadges.bitbadgeschain.badges.MsgRegisterAddresses", MsgRegisterAddresses],
    ["/bitbadges.bitbadgeschain.badges.MsgRequestTransferManager", MsgRequestTransferManager],
    ["/bitbadges.bitbadgeschain.badges.MsgSetApproval", MsgSetApproval],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateBytes", MsgUpdateBytes],
    ["/bitbadges.bitbadgeschain.badges.MsgNewBadge", MsgNewBadge],
    ["/bitbadges.bitbadgeschain.badges.MsgNewSubBadge", MsgNewSubBadge],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdatePermissions", MsgUpdatePermissions],
    ["/bitbadges.bitbadgeschain.badges.MsgFreezeAddress", MsgFreezeAddress],
    ["/bitbadges.bitbadgeschain.badges.MsgTransferManager", MsgTransferManager],
    ["/bitbadges.bitbadgeschain.badges.MsgRevokeBadge", MsgRevokeBadge],
    ["/bitbadges.bitbadgeschain.badges.MsgRequestTransferBadge", MsgRequestTransferBadge],
    ["/bitbadges.bitbadgeschain.badges.MsgSelfDestructBadge", MsgSelfDestructBadge],
    ["/bitbadges.bitbadgeschain.badges.MsgHandlePendingTransfer", MsgHandlePendingTransfer],
    ["/bitbadges.bitbadgeschain.badges.MsgPruneBalances", MsgPruneBalances],
    
];

export { msgTypes }