import React from "react";
import { RVerPromiseClient } from "./pb/rvapi/rvapi_grpc_web_pb";

export const ClientContext = React.createContext(new RVerPromiseClient(""));
export const WindowBaseURLContext = React.createContext("");
export const CheckedInContext = React.createContext(false);
