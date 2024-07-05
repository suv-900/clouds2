import React, { useState } from "react";

export default function EditUser(props:{
    enable:boolean,
    token:string,
    userAbout:string,
    username:string
}){
    const {enable,token,userAbout,username} = props
    const[about,setAbout] = useState("")
    const[errorMessage,setErrorMessage] = useState("");
    const[render,setRender] = useState(enable);

    function displayError(s:string){
        setErrorMessage(s);
        setTimeout(()=>{
            setErrorMessage("");
        },2000)
    }
    function close(){
        setTimeout(()=>{
            setRender(false)
        },800)
    } 
    async function update(){
        const response = await fetch(`http://localhost:8000/users/update-about`,{
            method:"PUT",
            headers:{"Authorization":token},
            body:about
        })
        if(response.status === 200){
            close()
        }else if(response.status === 400){
            displayError("please login")
        }else if(response.status === 500){
            displayError("server error.try again later")
        }
    }
    return(
        <div className={`useredit-container ${render?'slide-down':'slide-up'}`}>
            <div> 
                <div>Profile Information</div>
                <h2>{username}</h2>
                <label>About</label>
                <textarea
                className="useredit-abouttextarea" 
                onChange={(e)=>{
                        setAbout(e.target.value)
                    }}
                />
                <button onClick={()=>{
                    if(about === ""){
                        displayError("about cannot be empty")
                    }else{
                        update()
                    }
                }}>save</button>
                <p className="useredit-errormessage">{errorMessage}</p>
                </div>    
        </div>    
    )
}