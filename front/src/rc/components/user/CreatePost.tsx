import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

export default function CreatePost(){
    const[token,setToken] = useState("");
    const[title,setTitle] = useState("");
    const[body,setBody] = useState("");
    const[displayError,setDisplayError] = useState(false);
    const[error,setError] = useState("");
    
    const navigator = useNavigate();
    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token === null){
            navigator("/login")
        }else{
            setToken(token);
        }
    },[])
    async function createPost(){
        if(title.length === 0 || body.length === 0 || token.length === 0){
            return;
        }
        const headers={
            "Authorization":token
        }
        const post = {"post_title":title,"post_content":body}
        const requestBody = JSON.stringify(post)

        const response = await fetch("http://localhost:8000/createpost",{
            method:"POST",
            body:requestBody,
            headers:headers
        })
        if(response.ok){
            const postid = await response.json();
            setTimeout(()=>{
                navigator(`/post/view?id=${postid}`)
            },1800)
        }else if(response.status === 500){
            renderError("server error try again.")
        }else if(response.status === 400 || 401){
            renderError("please login.");
            setTimeout(()=>{
                navigator("/login")
            },3000)
        }
    }

    function vanishErrorMessage(){
        setTimeout(()=>{
            setDisplayError(false);
            setError("");
        },3000)
    }
    function renderError(msg:string){
        setError(msg);
        setDisplayError(true);
        vanishErrorMessage();
    }
    return(
        <div className="create-post-container">
            <div className="create-post-form">
                <div className="create-post-title">Write a post</div>
                <label className="create-post-label">title</label>
                <textarea
                    autoFocus
                    className="title-input" 
                    onChange={(e)=>{setTitle(e.target.value)}}
                    placeholder="title..."
                    wrap="soft"
                />
                
                <label className="create-post-label">body</label>
                <textarea
                    className="body-input" 
                    onChange={(e)=>{setBody(e.target.value)}}
                    placeholder="Write your post..."
                />
                <button className="submit-button" onClick={()=>{createPost()}}>create</button>
                <div>{displayError?error:""}</div>
            </div>
        </div>
    )
}