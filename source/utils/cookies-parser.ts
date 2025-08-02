export const getCookies = () => {
  const header_cookies = "";
  const youtubeCookies: { name: string; value: string }[] = [];

  for (const cookies of header_cookies.split(";")) {
    const [name, value] = cookies.trim().split("=");
    youtubeCookies.push({ name, value });
  }
  return youtubeCookies;
};
