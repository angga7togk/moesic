import ytdl from "@distube/ytdl-core";
import { getCookies } from "./cookies-parser.js";

const getYoutube = async (url: string) => {
  const uri = decodeURIComponent(url);

  const agent = ytdl.createAgent(getCookies());
  const stream = ytdl(uri, {
    agent,
    quality: "highestaudio",
    filter: "audioonly",
    highWaterMark: 1 << 25,
  });
  return stream;
};

export default {
  getYoutube,
};
