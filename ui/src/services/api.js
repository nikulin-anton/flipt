import axios from "axios";

export const Api = axios.create({
  baseURL: "/feature-flags/api/v1/",
});
