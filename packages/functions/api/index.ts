import { NowRequest, NowResponse } from "@vercel/node"
export default (req: NowRequest, res: NowResponse) => {
    res.send("Hello, unwanted guest!!")
}