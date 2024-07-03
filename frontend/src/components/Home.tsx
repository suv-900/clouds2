import React, { useEffect, useState } from "react";
import Post from "../types/Post";
import PostCule from "./posts/PostCule";
import Loading from "./Loading";

export default function Home(){
    const[posts,setPosts] = useState<Post[]>([]); 
    const[offset,setOffset] = useState(0);
    const[loading,setLoading] = useState(false);
    const[fetching,setFetching] = useState(false);

    useEffect(()=>{
        window.addEventListener("scroll",handleScroll)    
        startLoading()
        getPosts()
        
        return ()=>window.removeEventListener("scroll",handleScroll)
    },[])

    function handleScroll(){
        if(window.innerHeight + document.documentElement.scrollTop === document.documentElement.offsetHeight){
            setFetching(true);
        }
    }

    useEffect(()=>{
        if(!fetching) return;
        setOffset(prevValue=>prevValue+1)
        getPosts()
        setTimeout(()=>{
            setFetching(false)
        },3000)
    },[fetching])

    function startLoading(){
        setLoading(true);

        setTimeout(()=>{
            setLoading(false)
        },3000)
    }
    async function getPosts(){
        const response = await fetch(`http://localhost:8000/posts/getposts?offset=${offset}&limit=5`)
        const res = await response.json();
        if(res === null) return
        for(let i=0;i<res.length;i++){
            const k = res[i];
            const post = new Post(
                    k.Post_id,
                    k.Post_title,
                    "",
                    k.Author_name,
                    0,
                    k.Post_likes,
                    false,
                    false,
                    ""
            )
            setPosts(prevArr=>[...prevArr,post])
        }
    }

    return(
        <div>
            <Loading enable={loading}/>
            {!loading?
            <div>
                <h3>Home</h3>
                {posts.map(post=><PostCule post={post}/>)}

                {fetching && <span className="scroll-loader"></span>} 
            </div>:<></>}
        </div>
    )
}