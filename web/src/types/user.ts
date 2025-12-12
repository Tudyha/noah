export type LoginRequest = {
  login_type: number;
  username: string;
  password?: string;
  code?: string;
};

export type LoginResponse = { token: string };

export type UserResponse = {
  id: string;
  nickname: string;
  avatar: string;
  work_space_list: WorkSpaceResponse[];
};

export type WorkSpaceResponse = {
  id: string;
  name: string;
  description: string;
  app_list: WorkSpaceAppResponse[];
};

export type WorkSpaceAppResponse = {
  id: string;
  name: string;
  description: string;
};
