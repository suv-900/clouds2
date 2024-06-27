import React, { useState,useEffect } from 'react'
import {useNavigate} from 'react-router-dom'
import { checkEmail, checkPassword } from '../../utils';
import { FormError } from '../errors/FormErrors';
import Loading from '../Loading';

//checkusername sends 2 requests for bigger names
export default function Register(){
    const navigate = useNavigate();
    
    const[username,setUsername] = useState("");
    const[password,setPassword] = useState("");
    const[email,setEmail] = useState("");
   
    const[usernameError,setUsernameError] = useState(false);
    const[usernameErrorMessage,setUsernameErrorMessage] = useState<string>();
    const[emailError,setEmailError] = useState(false);
    const[emailErrorMessage,setEmailErrorMessage] = useState<string>();
    const[passwordError,setPasswordError] = useState(false);
    const[passwordErrorMessage,setPasswordErrorMessage] = useState<string>();

    const[usernameValid,setUsernameValid] = useState(false);
    const[emailValid,setEmailValid] = useState(false);
    const[passwordValid,setPasswordValid] = useState(false);

    const[error,setError] = useState(false);
    const[errorMessage,setErrorMessage] = useState<string>("");
    const[loading,setLoading] = useState(false);

    let timer:NodeJS.Timeout;
    function delayRequest(username:string){
        console.log(username);
        if(timer != null){
            console.log("clearing timer");
            clearTimeout(timer);
        }

        timer = setTimeout(()=>{
            checkUsernameExists(username);
        },4000)

    }
    async function checkUsernameExists(username:string){
        console.log(username);
        const response = await fetch("http://localhost:8000/checkusername",{
            body:JSON.stringify(username),
            method:"POST"
        })
        if(response.status === 200){
            setUsername(username);
            setUsernameValid(true); 
        }else if(response.status === 409){
            setUsernameError(true);
            setUsernameErrorMessage("username exists.");
            //refactor
            setUsernameValid(false);
        }else if(response.status === 500){
            displayError("server error.try again");
        }else{
            displayError("server error.try again");
        }
    }

    function displayError(msg:string){
        setErrorMessage(msg);
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
    async function register(){
        if(!usernameValid || !passwordValid || !emailValid){
            displayError("invalid form");
            return;
        }

        const requestbody = JSON.stringify({username,email,password});
        const response = await fetch("http://localhost:8000/register",{
            body:requestbody,
            method:"POST"
        })

        if(response.status === 200){
            const token = await response.json();
            saveToken(token);
        }else if(response.status === 500){
            displayError("server error.try again");
        }else if(response.status === 409){
            //conflict
        }else{
            displayError("server error.try again");
        }
    }

    return(
        <div> 
        {!loading?
        <div>
            <label>username</label><br></br>
            <input type="text" 
            onChange={(e)=>{
                if(e.target.value.length === 0){
                    setUsernameError(true);
                    setUsernameErrorMessage("username cannot be blank");
                    return;
                }
                delayRequest(e.target.value)
                setUsernameError(false);
            }}
            />
            <br></br>
            <FormError render={usernameError} msg={usernameErrorMessage}/>
            <br></br>
            
            <label>password</label><br></br>
            <input type="password"
            onChange={(e)=>{
                // if(!checkPassword(e.target.value)){
                //     setPasswordValid(false);
                // }else{
                //     setPassword(e.target.value);
                //     setPasswordValid(true);
                // }
                if(e.target.value.length === 0){
                    setPasswordValid(false);
                    setPasswordError(true);
                    setPasswordErrorMessage("password cannot be blank");
                    return;
                }
                setPasswordValid(true);
                setPasswordError(false);
                setPassword(e.target.value);
            }} 
            />
            <br></br>
            <FormError render={passwordError} msg={passwordErrorMessage}/>
            <br></br>

            <label>email</label><br></br>
            <input type="email"
            onChange={(e)=>{
                if(e.target.value.length === 0){
                    setEmailError(true);
                    setEmailErrorMessage("email cannot be blank");
                    return;
                }
                if(!checkEmail(e.target.value)){
                    setEmailError(true);
                    setEmailErrorMessage("invalid email");
                }else{
                    setEmailError(false);
                    setEmail(e.target.value);
                    setEmailValid(true);
                }
            }} 
            />
            <br></br>
            <FormError render={emailError} msg={emailErrorMessage}/>
            <br></br>
            
            <button onClick={()=>{register()}}
            >register</button>
            
            {error && errorMessage ?<div>{errorMessage}</div>:<></>}
        </div>
        :<></>}
            <Loading enable={loading}/>
        </div>
    )
}


