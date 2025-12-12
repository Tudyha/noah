export * from "./user";
export * from "./common";
export * from "./client";

export type BaseResponse<T> = {
  code: number;
  msg: string;
  data: T | null;
};

export type PageResponse<T> = {
  total: number;
  list: T[];
};