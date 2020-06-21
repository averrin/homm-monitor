function fillTowns(report, player) {
    player.querySelector("#townsCount").innerText = `${report.townsCount}`;

    player.querySelector("#gold").innerText = `${report.resources.gold} (${report.resources.goldSpent})`;
    player.querySelector("#wood").innerText = `${report.resources.wood} (${report.resources.woodSpent})`;
    player.querySelector("#mercury").innerText = `${report.resources.mercury} (${report.resources.mercurySpent})`;
    player.querySelector("#ore").innerText = `${report.resources.ore} (${report.resources.oreSpent})`;
    player.querySelector("#sulfur").innerText = `${report.resources.sulfur} (${report.resources.sulfurSpent})`;
    player.querySelector("#crystal").innerText = `${report.resources.crystal} (${report.resources.crystalSpent})`;
    player.querySelector("#gems").innerText = `${report.resources.gems} (${report.resources.gemsSpent})`;


    let towns = player.querySelector("#towns");

    towns.innerHTML = '';
    report.towns.forEach(town => {
        let item = document.createElement('div');
        item.className = "row townLine";
        let ava = document.createElement('div');
        ava.className = "townImage";
        ava.style.backgroundImage = `url(/towns/${town.type}.gif)`;
        if (!town.hasFort) {
            ava.style.backgroundImage = `url(/towns/${town.type}_0.png)`;
        }
        ava.style.borderColor = colors[report.color];
        item.append(ava);

        let townInfo = document.createElement('div');
        townInfo.className = "townInfo row";

        let c1 = document.createElement('div');
        c1.className = "column";
        let c2 = document.createElement('div');
        c2.className = "column";

        let name = document.createElement('div');
        name.className = "townName value";
        name.innerText = town.name;
       
        let army = document.createElement('div');
        army.innerHTML = `Army: <span class='value'>${town.guardsValue}</span>&nbsp&nbsp`;

        c1.append(name);
        c2.append(army);

        //let lables = ["I", "II", "III", "IV", "V"];
        for (let index = 0; index < town.spells.length - 2; index++) {
            let spells = document.createElement('div');
            spells.className = "row spellsRow";
            //spells.innerHTML = `${lables[index]}&nbsp;`;
            for (let n = 0; n <town.spells[index].length - 1; n++) {
                const spell = town.spells[index][n];
                if (spell > 69) continue;
                spells.innerHTML += `<img class="spellIcon" src="/spells/${spell}.png">`;
            }
            if (index == 3 && town.spells.length == 5) {
                //spells.innerHTML = `&nbsp;&nbsp;${lables[index]}&nbsp;`;
                for (let n = 0; n <town.spells[4].length - 1; n++) {
                    const spell = town.spells[4][n];
                    if (spell > 69) continue;
                    spells.innerHTML += `<img class="spellIcon" src="/spells/${spell}.png">`;
                }
            }
            [c1, c2][index%2].append(spells);      
        }

        townInfo.append(c1);
        townInfo.append(c2);
        
        item.append(townInfo);
        towns.append(item);
    });
}