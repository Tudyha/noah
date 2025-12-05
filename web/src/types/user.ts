export type LoginRequest = {
  loginType: number;
  username: string;
  password?: string;
  code?: string;
};
export type LoginResponse = { token: string };
export type UserResponse = {
  id: string;
  nickname: string;
  avatar: string;
  workSpaceList: WorkSpaceResponse[];
};

export type WorkSpaceResponse = {
  id: string;
  name: string;
  description: string;
  appList: WorkSpaceAppResponse[];
};

export type WorkSpaceAppResponse = {
  id: string;
  name: string;
  description: string;
};
