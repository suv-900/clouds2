import React, { useEffect, useState } from "react";
import Post from "../types/Post";
import PostCule from "./posts/PostCule";
import Loading from "./Loading";

export default function Home(){
    const[posts,setPosts] = useState<Post[]>([]); 
    const[offset,setOffset] = useState(0);
    const[render,setRender] = useState(false);
    const[loading,setLoading] = useState(false);
    const[isFetching,setisFetching] = useState(false);

    useEffect(()=>{
        window.addEventListener("scroll",handleScroll)    
        startLoading()
        getPosts()
        
        return ()=>window.removeEventListener("scroll",handleScroll)
    },[])

    function handleScroll(){
        if(window.innerHeight + document.documentElement.scrollTop === document.documentElement.offsetHeight){
            setOffset(prevValue=>prevValue+1)
            getPosts()
        }
    }
    function startLoading(){
        setLoading(true);

        setTimeout(()=>{
            setLoading(false)
        },3000)
    }
    async function getPosts(){
        setisFetching(true)
        const response = await fetch(`http://localhost:8000/posts/getposts?offset=${offset}&limit=5`)
        const res = await response.json();
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
            // setPosts(prevArr=>[...prevArr,post])
            posts.push(post)
        }
        setTimeout(()=>{
            setisFetching(false)
        },5000)
    }

    return(
        <div>
            <Loading enable={loading}/>
            {!loading?
            <div>
                <h3>Home</h3>
                {posts.map(post=><PostCule post={post}/>)}

                {isFetching && <span className="scroll-loader"></span>} 
            </div>:<></>}
        </div>
    )
}