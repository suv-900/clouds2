import React, { useState,useEffect } from 'react'
import {useNavigate} from 'react-router-dom'
import MainPage from '../Home';
import Loading from '../Loading';

//disable button first render
export default function Login(){
    const navigate = useNavigate();
    //function returns the copy of the var and function to update the copy
    const[username,setUsername] = useState<string>("");
    const[password,setPassword] = useState<string>("");
    const[loginSuccess,setLoginSuccess] = useState(false);
    const[error,setError] = useState<boolean>();
    const[errorMessage,setErrorMessage] = useState<string>();
    const[loading,setLoading] = useState(false);
    
    // const navigator = useNavigate();
    useEffect(()=>{
        if(loginSuccess){
            // navigator("/home");
        }
    },[loginSuccess])

    
    function displayError(message:string){
        setErrorMessage(message);
        setError(true);

        setTimeout(()=>{
            setError(false);
            setErrorMessage("");
        },3000)
    }
    function saveToken(token:string){
        setLoading(true);
        localStorage.setItem("token",token);
        setTimeout(()=>{
            setLoading(false);
            navigate("/home");
        },4000)
    }
    async function login(){
        if(username.length === 0 || password.length === 0){
            displayError("invalid form");
            return;
        }
    
        const requestBody= {username,password}
        const body = JSON.stringify(requestBody);

        const response = await fetch("http://localhost:8000/login",{
            body:body,
            method:"POST"
        })

        if(response.ok){
            const token = await response.json()
            saveToken(token);
        }else if(response.status === 401){
            displayError("wrong password");
        }else if(response.status === 404){
            displayError("user not found");
        }else if(response.status === 500){
            displayError("server error.try again!");
        }else{
            displayError("server error.try again!");
        }
    } 

    return(
        <div >
        {loading?<></>:   
        <div className="login-form">
            <label>username</label><br></br>
            <input type="text" onChange={(e)=>{setUsername(e.target.value)}} /><br></br>
            
            <label>password</label><br></br>
            <input type="text" onChange={(e)=>{setPassword(e.target.value)}} /><br></br>
            
            <button onClick={()=>{login()}}>login</button>
            {error && errorMessage ?<div>{errorMessage}</div>:<></>}
        </div>
        }   
        <Loading enable={loading}/>
        </div>
    )
}

