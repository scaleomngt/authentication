// The 'card.aleo' program.
program card.aleo {
    struct CardInfo {
    	gender: field,
    	name: field,
    	age: u32,
        nation: field,
        birthdate: field,
        addr: address,
    }

    record Card {
        owner: address,
        gates: u64,
        id: field,
    	name: field,
    	gender: field,
        nation: field,
        birthdate: field,
        age: u32,
    }
    
    record CardResult {
        owner: address,
        gates: u64,
        id: field,
        name: field,
        gender: field,
        result: u32,
    }

    transition submit(public info: CardInfo) -> Card {
        let id: field = BHP256::hash_to_field(info.addr);
        let result: u64 = 0u64;
        return Card {
            owner: self.caller,
            gates: 0u64,
            id,
            age: info.age,
            name: info.name,
            gender: info.gender,
            nation: info.nation,
            birthdate: info.birthdate,
        };
    }


    transition calculate_age(public card: Card) -> CardResult {
    	let result: u32 = 0u32;
    	let b: u32 = card.age;

        if(b < 18u32){
            result = 1u32;
        }else {
            result =  0u32;
        }
        
        return CardResult {
         owner: card.owner,
         gates: 0u64,
         id: card.id,
         name: card.name,
         gender: card.gender,
         result: result,
       };
    }
}
