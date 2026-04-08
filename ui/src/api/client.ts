import ky from "ky";

const TOKEN_KEY = "agenthound_token";

export const api = ky.create({
  prefixUrl: "/api/v1",
  hooks: {
    beforeRequest: [
      (request) => {
        const token = localStorage.getItem(TOKEN_KEY);
        if (token) {
          request.headers.set("Authorization", `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      (_request, _options, response) => {
        if (response.status === 401) {
          localStorage.removeItem(TOKEN_KEY);
          localStorage.removeItem("agenthound_user");
          if (window.location.pathname !== "/login") {
            window.location.href = "/login";
          }
        }
      },
    ],
  },
});
