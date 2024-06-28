class node<T>{
    value:T;
    prev:node<T> | null;
    next:node<T> | null;

    constructor(v:T){
        this.value = v;
        this.prev = null;
        this.next = null;
    }
}

export class LinkedList<T>{
   head:node<T> | null;
   tail:node<T> | null;
   
   constructor(){
    this.head = this.tail = null;
   }

   push(v:T){
    let newNode = new node(v);
    
    if(this.head == null && this.tail == null){
        this.head = this.tail = newNode;
    }else if(this.tail != null){
        this.tail.next = newNode;
        newNode.prev = this.tail;
        this.tail = newNode;
    }
   }

   pop(){
    if(this.tail != null){
        let temp = this.tail;
        this.tail = this.tail.prev;
        if(this.tail){
            this.tail.next = null;
            temp.prev = null;
        }
    }
   }
    
    push_front(v:T){
        let newNode = new node(v);

        if(this.head === null && this.tail === null){
            this.head = this.tail = newNode;
        }else if(this.head !== null){
            newNode.next = this.head;
            this.head.prev = newNode;
            this.head = newNode;
        }
    }

    remove(v:T){
        let temp = this.head;
        if(temp === null) return;
        
        while(temp.next !== null){
           if(temp.value === v){
                if(temp.prev === null && this.head !== null){
                    this.head = this.head.next;
                    if(this.head != null){
                        this.head.prev = null;
                    }
                }else if(temp.next === null && this.tail != null){
                    if(this.tail.prev != null){
                        this.tail = this.tail.prev;
                    }
                }else if(temp.prev != null && temp.next != null){
                    temp.prev.next = temp.next;
                    temp.next.prev = temp.prev;
                }
           }
           temp = temp.next; 
        }
    }

    contains(v:T):boolean{
        let temp = this.head;
        if(temp === null) return false;

        while(temp != null){
            if(temp.value === v) return true;
            temp = temp.next;
        }
        return false;
    }
}