import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import Post from "../types/Post";
import PostCule from "../posts/PostCule";

export default function Home(){
    const[token,setToken] = useState<string>();
    const[posts,setPosts] = useState<Post[]>([]); 
    const[offset,setOffset] = useState(0);
    const[render,setRender] = useState(false);
    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token !== null){
            setToken(token);
        }
        getPosts()
    },[])
    // const location = useLocation();
    // const t = location.state || "";
    // setToken(t);
    async function getPosts(){
        const response = await fetch(`http://localhost:8000/getposts/${offset}`)
        const res = await response.json();

        for(let i=0;i<res.length;i++){
            const k = res[i];
            const post = new Post(
                    k.Post_id,
                    k.Post_title,
                    k.Post_content,
                    k.Author_name,
                    k.Author_id,
                    k.Post_likes,
                    k.Createdat
            ) 
            posts.push(post);    
        }
        setRender(true)
    }

    return(
        <div>
            <h3>Home</h3>
            {render?posts.map(post=><PostCule post={post}/>):<></>}
        </div>
    )
}