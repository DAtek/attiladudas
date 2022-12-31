import type {FieldError} from "@/utils/api_client";


export function getErrorMessage(err: FieldError): string {
  const errorStr = errorTypeToHumanReadable[err.type]
  if (!errorStr) {
    throw `Error string isn't defined for '${err.type}'`
  }

  return errorStr
}

export function getErrorsForField(field: string, errors: FieldError[]): FieldError[] {
  return errors.filter((item) => item.location === field)
}

export const errorTypeToHumanReadable = {
    REQUIRED: "The field is required",
    INVALID_ISO8601: "The date is not in ISO8601 format",
    ALREADY_EXISTS: "Already exists",
    PLAYER_WITH_THIS_NAME_ALREADY_JOINED: "Player with this name already joined",
    ROOM_IS_FULL: "This room already has 2 players",
    BOTH_PLAYERS_MUST_JOIN: "Both players must join before you can pick a side",
    SIDE_ALREADY_TAKEN: "This side is already taken",
}


export type errorMappingKey = keyof typeof errorTypeToHumanReadable
