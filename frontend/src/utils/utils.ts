function checkPassword(password:string):boolean{
    if(password.length === 0) return false;

    const len =(password.length >=8 && password.length <= 20)?true:false;
    const digit = /\d/.test(password);
    const specialChar = /[@#$%&]/.test(password);
    
    return len && digit && specialChar;
}
function checkEmail(email:string):boolean{
    if(email.length === 0) return false;

    return /@gmail.com/.test(email);
}

// function decodeTime(s:string):string{
//     const currentTime = Date() 
//     let date="";
//     let time="";
//     let i = 0;
//     while(s[i] !== 'T'){
//         date += s[i]
//     }
// }

export {checkPassword,checkEmail};