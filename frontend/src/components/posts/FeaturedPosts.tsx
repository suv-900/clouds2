import React, { useEffect, useState } from "react";
import { PostMetaData } from "../../types/PostMetaData";

function FeaturedPosts(){
    const[offset,setOffset] = useState(0)
    const[loading,setLoading] = useState(true)
    const[list,setList] = useState<PostMetaData[]>([])
    useEffect(()=>{
        // getData()
    },[])
    
    useEffect(()=>{
        setLoading(true)
        // getData()
        setTimeout(()=>{
            setLoading(false)
        },2000)
    },[offset])

    async function getData(){
        const response = await fetch(`http://localhost:8000/posts/get-featured-posts/${offset}`)
        const res = await response.json()

        let dataList:PostMetaData[] = [];
        for(let i = 0;i<res.length;i++){
            const k = res[i];
            const data = new PostMetaData(k.id,k.title)
            dataList.push(data)
        }
        setList(dataList)
    }
    return(
        <div>
            <ul>
            {list.map((post)=>
            <li>
                <a href={`http://localhost:3000/post/view?id=${post.id}`}>{post.title}</a>
            </li>
        
            )}
            </ul>
            <button onClick={()=>{setOffset(offset+1)}}>load more...</button>
        </div>
    )
}
