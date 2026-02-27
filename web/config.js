const CONFIG = {
  BASE_URL: "http://127.0.0.1:8001",
  ENDPOINT: "/book",
};

function getUrl(id = null) {
  let url = CONFIG.BASE_URL + CONFIG.ENDPOINT;
  if (id) {
    url += "/" + id;
  }
  return url;
}
