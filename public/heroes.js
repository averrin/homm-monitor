function fillHeroes(report, player) {
    player.querySelector("#heroesCount").innerText = `${report.heroesCount}`;
    player.querySelector("#mpSpent").innerText = `${report.totalMPSpent}`;
    player.querySelector("#expEarned").innerText = `${report.totalExp}`;

    let heroes = player.querySelector("#heroes");

    heroes.innerHTML = '';
    report.heroes.forEach(hero => {
        let item = document.createElement('div');
        item.className = "row";
        let ava = document.createElement('div');
        ava.className = "heroImage";
        ava.style.backgroundImage = `url(/heroes/${hero.id}.png)`;
        ava.style.borderColor = colors[report.color];
        ava.innerHTML = `<img class="genderSymbol" src="/heroes/male.svg">`;
        if (hero.gender == 1) {
            ava.innerHTML = `<img class="genderSymbol" src="/heroes/female.svg">`;
        }
        item.append(ava);
        let heroInfo = document.createElement('div');
        heroInfo.className = "heroInfo row";

        let c1 = document.createElement('div');
        c1.className = "column";
        let c2 = document.createElement('div');
        c2.className = "column";

        let name = document.createElement('div');
        name.className = "heroName value";
        name.innerText = hero.name;
        if (hero.inGarrison) {
            name.innerHTML += `<img class="skillIcon" src="/towns/castle.gif" style="margin-left: 6px;">`;
        }
        let moves = document.createElement('div');
        moves.innerHTML = `MP: <span class='value'>${hero.movement}/${hero.maxMovement}</span>&nbsp&nbsp`;
        if (hero.sec_skills[2] > 0) {
            moves.innerHTML += `<img class="skillIcon" src="/skills/${2}_${hero.sec_skills[2]}.png">`;
        }
        let exp = document.createElement('div');
        exp.innerHTML = `XP: <span class='value'>${hero.level} (${hero.experience})</span>`;

        let army = document.createElement('div');
        army.innerHTML = `Army: <span class='value'>${hero.armyPower} (${hero.armyValue})</span>&nbsp&nbsp`;
        for (let i = 22; i <= 23; i++) {
            if (hero.sec_skills[i] > 0) {
                army.innerHTML += `<img class="skillIcon" src="/skills/${i}_${hero.sec_skills[i]}.png">`;
            }
        }

        let stats = document.createElement('div');
        stats.innerHTML = `Stats: <span class='value'>${hero.attack}</span>-`;
        stats.innerHTML += `<span class='value'>${hero.defense}</span>-`;
        stats.innerHTML += `<span class='value'>${hero.spell_power}</span>-`;
        stats.innerHTML += `<span class='value'>${hero.knowledge}</span>`;
        //stats.innerHTML += `&nbsp[<span class='value'>${hero.moraleBonus} / ${hero.luckBonus}</span>]`;
        
        let mana = document.createElement('div');
        mana.innerHTML = `Mana: <span class='value'>${hero.mana}</span>&nbsp&nbsp`;

        for (let i = 15; i <= 17; i++) {
            if (hero.sec_skills[i] > 0) {
                mana.innerHTML += `<img class="skillIcon" src="/skills/${i}_${hero.sec_skills[i]}.png">`;
            }
        }


        c1.append(name);
        c1.append(moves);
        c1.append(exp);
        c2.append(army);
        c2.append(stats);
        c2.append(mana);
        heroInfo.append(c1);
        heroInfo.append(c2);
        
        item.append(heroInfo);
        heroes.append(item);
    });
}