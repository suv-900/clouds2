import React, { useEffect, useState } from "react";
import Post from "../types/Post";
import { useNavigate } from "react-router-dom";

export default function PostCule(props:{
    post:Post
}){
    const[post,setPost] = useState(props.post);
    const[likes,setLikes] = useState(post.likes);
    const[token,setToken] = useState<string>();
    
    const authorLink = `http://localhost:3000/user/${post.authorid}`;
    const postLink = `http://localhost:3000/post/view?id=${post.id}`;

    const navigate = useNavigate();
   
    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token !== null){
            setToken(token);
        }
    },[])


    async function likePost(){
        if(!token){
            navigate("/login")
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/likepost/${post.id}`,{
            method:"POST",
            headers:headers
        })
        if(response.ok){
            setLikes(likes+1);
        }else if(response.status === 400){
            navigate("/login")
        }else if(response.status === 401){
            navigate("/register")
        }
    }
    
    async function dislikePost(){
        if(!token){
            navigate("/login")
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/likepost/${post.id}`,{
            method:"POST",
            headers:headers
        })
        if(response.ok){
            setLikes(likes-1);
        }else if(response.status === 400){
            navigate("/login")
        }else if(response.status === 401){
            navigate("/register")
        }
    }
    return(
        <div>
            <a href={postLink}>{post.title}</a>
            <a href={authorLink} >{post.authorname}</a>
            <div>{post.createdat}</div>
            <div>{likes}</div>
            <button onClick={()=>{likePost()}}>like</button>
            <button onClick={()=>{dislikePost()}}>dislike</button>
        </div>
    )
}