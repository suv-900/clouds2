import React from "react";
import Post from "../types/Post";
import TopHeader from "../components/TopHeader";
import PostContent from "./PostContent";
import PostComments from "./PostComments";
import Comment from "../types/Comment";

export default function PostComponent(props:{
    post:Post | undefined,
    comments:Comment[],
    token:string | null,
}){
    return(
        <div>
            <TopHeader/>
            <PostContent 
                token={props.token} 
                postContent={props.post}               
            />
            <PostComments
                token={props.token} 
                comments={props.comments}
                postid={props.post?props.post.id:undefined}
            />

        </div>
    )
}