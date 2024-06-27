import React, { useEffect, useState } from "react";
import Post from "../types/Post";
import "../css/styles.css";
import { useNavigate } from "react-router-dom";
import CustomError from "../components/CustomError";
import PostFooter from "./PostFooter";

export default function PostContent(props:{
    tokenFound:boolean,
    postContent:Post | undefined,
}){
    const post = props.postContent;  
    return(
        <div>
            {post ?
            <div>
            <div className="post">
                <img src={process.env.PUBLIC_URL+'/retrobuttons/book.gif'} alt="image not found"/>  
                <div className="post-header">
                    <h2 className="post-title">{post.title}</h2>
                    <h5 className="post-authorname">{post.authorname}</h5> 
                </div>
                <p className="post-content">{post.content}</p>
                <div className="post-footer">
                   <PostFooter tokenFound={props.tokenFound} postid={post.id} likes={post.likes} /> 
                </div>
            </div>
            <div>
                <img src={process.env.PUBLIC_URL+'/retrobuttons/linkexchange.gif'} alt="image not found"/>  
                <img src={process.env.PUBLIC_URL+'/retrobuttons/computer-2.jpg'} alt="image not found"/>  
                <img src={process.env.PUBLIC_URL+'/retrobuttons/gaygoogle.gif'} alt="image not found" />  
                <img src={process.env.PUBLIC_URL+'/retrobuttons/iexplorer.gif'} alt="image not found" />  
                <img src={process.env.PUBLIC_URL+'/retrobuttons/win95.gif'} alt="image not found" />  
            </div>
            </div>
            :<></>}
        </div>
    )
}