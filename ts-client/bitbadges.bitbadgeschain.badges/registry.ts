import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateAddressMappings } from "./types/badges/tx";
import { MsgUpdateCollection } from "./types/badges/tx";
import { MsgDeleteCollection } from "./types/badges/tx";
import { MsgUpdateUserApprovedTransfers } from "./types/badges/tx";
import { MsgTransferBadges } from "./types/badges/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/bitbadges.bitbadgeschain.badges.MsgCreateAddressMappings", MsgCreateAddressMappings],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateCollection", MsgUpdateCollection],
    ["/bitbadges.bitbadgeschain.badges.MsgDeleteCollection", MsgDeleteCollection],
    ["/bitbadges.bitbadgeschain.badges.MsgUpdateUserApprovedTransfers", MsgUpdateUserApprovedTransfers],
    ["/bitbadges.bitbadgeschain.badges.MsgTransferBadges", MsgTransferBadges],
    
];

export { msgTypes }