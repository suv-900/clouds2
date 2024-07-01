import React, { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import Post from "../types/Post";
import PostCule from "./posts/PostCule";
import Loading from "./Loading";

export default function Home(){
    const[token,setToken] = useState<string>();
    const[posts,setPosts] = useState<Post[]>([]); 
    const[offset,setOffset] = useState(0);
    const[render,setRender] = useState(false);
    const[loading,setLoading] = useState(false);
    
    useEffect(()=>{
        startLoading()
        const token = localStorage.getItem("token");
        if(token !== null){
            setToken(token);
        }
        getPosts()
    },[])
    // const location = useLocation();
    // const t = location.state || "";
    // setToken(t);

    function startLoading(){
        setLoading(true);

        setTimeout(()=>{
            setLoading(false)
        },3000)
    }
    async function getPosts(){
        const response = await fetch(`http://localhost:8000/posts/getall`)
        const res = await response.json();
        console.log(res);
        for(let i=0;i<res.length;i++){
            const k = res[i];
            const post = new Post(
                    k.Post_id,
                    k.Post_title,
                    k.Post_content,
                    k.Author_name,
                    k.Author_id,
                    k.Post_likes,
                    false,
                    false,
                    k.CreatedAt
            )
            console.log(post); 
            posts.push(post);    
        }
    }

    return(
        <div>
            <Loading enable={loading}/>
            {!loading?
            <div>
                <h3>Home</h3>
                {posts.map(post=><PostCule post={post}/>)}
            </div>:<></>}
        </div>
    )
}