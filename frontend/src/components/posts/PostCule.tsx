import React, { useEffect, useState } from "react";
import Post from "../../types/Post";
import { useNavigate } from "react-router-dom";
import {BiSolidLike,BiSolidDislike} from "react-icons/bi"

export default function PostCule(props:{
    post:Post
}){
    const[post,setPost] = useState(props.post);
    const[likes,setLikes] = useState(post.likes);
    const[token,setToken] = useState<string>();
    
    const navigate = useNavigate();
   
    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token !== null){
            setToken(token);
        }
    },[])

    function gotoPost(){
        navigate(`/post/view?id=${post.id}`)
    }

    function gotoUser(){
        navigate(`/users/${post.authorname}`)
    }

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
        <div onClick={()=>{gotoPost()}} className="post-container">
            <a onClick={()=>{gotoPost()}}  className="post-link">{post.title}</a>
            <a onClick={()=>{gotoUser()}} className="author-name" >{post.authorname} author</a>
            <div className="post-likes">{likes} likes</div>
            
            <BiSolidLike size="20" className="like-button"  onClick={()=>{likePost()}}>like</BiSolidLike>
            <BiSolidDislike size="20" className="dislike-button" onClick={()=>{dislikePost()}}>dislike</BiSolidDislike>
        </div>
        <div className="post-separator"></div>
        </div>
    )
}