import { useState,useEffect } from 'react'
import {useNavigate} from 'react-router-dom'
import { FormError } from '../errors/FormErrors';

//disable button first render
export default function Login(){
    //function returns the copy of the var and function to update the copy
    const[username,setUsername] = useState<string>("");
    const[password,setPassword] = useState<string>("");
    const[loginSuccess,setLoginSuccess] = useState(false);
    const[error,setError] = useState(false);
    const[errorMessage,setErrorMessage] = useState<string>();
    
    const navigator = useNavigate();
    useEffect(()=>{
        if(loginSuccess){
            navigator("/home");
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
    
    async function login(){
        if(username.length === 0 || password.length === 0){
            displayError("invalid form");
            return;
        }
    
        const requestBody= {username,password}
        const body = JSON.stringify(requestBody);

        const response = await fetch("http://localhost:8000/users/login",{
            body:body,
            method:"POST"
        })

        if(response.status === 200){
            const token = await response.json();
            localStorage.setItem("token",token);
            localStorage.setItem("username",username);
            localStorage.setItem("viewer-loggedin","true");
            setTimeout(()=>{
                setLoginSuccess(true)
            },2000)
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
        <div className="form-container">
            <div className="form-title">Login to your account</div>
            <label className="form-label">username</label>
            <input type="text" className="form-input" onChange={(e)=>{setUsername(e.target.value)}} /><br></br>
            
            <label className="form-label">password</label>
            <input type="text" className="form-input" onChange={(e)=>{setPassword(e.target.value)}} /><br></br>
            <FormError msg={errorMessage} enable={error}/>    
            <button className="form-button" onClick={()=>{login()}}>login</button>
        </div>
    )
}

