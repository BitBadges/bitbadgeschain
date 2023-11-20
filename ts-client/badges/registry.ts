import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgDeleteCollection } from "./types/badges/tx";
import { MsgCreateCollection } from "./types/badges/tx";
import { MsgCreateAddressMappings } from "./types/badges/tx";
import { MsgUniversalUpdateCollection } from "./types/badges/tx";
import { MsgTransferBadges } from "./types/badges/tx";
import { MsgUpdateUserApprovals } from "./types/badges/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/badges.MsgDeleteCollection", MsgDeleteCollection],
    ["/badges.MsgCreateCollection", MsgCreateCollection],
    ["/badges.MsgCreateAddressMappings", MsgCreateAddressMappings],
    ["/badges.MsgUniversalUpdateCollection", MsgUniversalUpdateCollection],
    ["/badges.MsgTransferBadges", MsgTransferBadges],
    ["/badges.MsgUpdateUserApprovals", MsgUpdateUserApprovals],
    
];

export { msgTypes }