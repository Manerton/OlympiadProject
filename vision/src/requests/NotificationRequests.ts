import axios from "axios";
import { NOTIFY } from "../config/api";

export async function axiosSendCode(email: string) {
    const res = await axios.post(
        NOTIFY.sendCode, {"email": email, "requestToken": "1234567890"},  // TODO!! 
    );
}