export * from "./user";
export * from "./common";

export type BaseResponse<T> = {
  code: number;
  msg: string;
  data: T | null;
};
