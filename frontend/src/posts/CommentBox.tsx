import React, { useEffect, useState } from "react";

export default function CommentBox(props:{
    postid:number | undefined
}){
    const[tokenFound,setTokenFound] = useState(false);
    const[token,setToken] = useState<string>();
    const[comment,setComment] = useState<string>();
    const[displayError,setDisplayError] = useState(false);

    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token !== null){
            setToken(token);
            setTokenFound(true);
        }
        console.log(props.postid === undefined)
    },[])
    
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
        {tokenFound && props.postid?
        <div>
            <label>add comment</label>
            <input type="text" onChange={(e)=>{setComment(e.target.value)}}/>
            <button onClick={()=>{addComment()}}>submit</button>
        </div>:<></>}
    </div>)
}