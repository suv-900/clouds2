import React from "react";
import Post from "../types/Post";
import TopHeader from "../components/TopHeader";
import PostContent from "./PostContent";
import CommentBox from "./CommentBox";
import PostComments from "./PostComments";
import Comment from "../types/Comment";


export default function PostComponent(props:{
    post:Post | undefined,
    comments:Comment[],
    tokenFound:boolean,
}){
    console.log("post component");
    return(
        <div>
            <TopHeader/>
            <PostContent 
                tokenFound={props.tokenFound} 
                postContent={props.post}               
            />
            <CommentBox postid={props.post?props.post.id:undefined} />
            <PostComments 
                comments={props.comments}
            />

        </div>
    )
}