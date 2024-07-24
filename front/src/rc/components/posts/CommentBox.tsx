import React, { useContext, useEffect, useState } from "react";
import "../../css/styles.css"
import { AuthContext } from "./PostViewer";

export default function CommentBox(props:{
    postid:number | undefined
}){

    const token = useContext(AuthContext)
    const[comment,setComment] = useState<string>();
    const[displayError,setDisplayError] = useState(false);

    function vanishErrorMessage(){
        setTimeout(()=>{
            setDisplayError(false);
        },3000)
    }
    function renderError(){
        setDisplayError(true);
        vanishErrorMessage();
    }
    async function addComment(){
        if(comment === undefined || comment.length === 0){
            renderError();
        }
    }
    return(
    <div>
        {token.length !== 0 && props.postid?
        <div className="comment-box">
            <label className="post-content">Add comment:</label><br></br>
            <textarea
                onChange={(e)=>{setComment(e.target.value)}}
                /><br></br>
            <button className="post-like-button"onClick={()=>{addComment()}}>submit</button>
        </div>:<></>}
    </div>)
}