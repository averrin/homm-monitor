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

        let lc = document.createElement('div');
        lc.className = "column";
        item.append(lc);

        let ava = document.createElement('div');
        ava.className = "townImage";
        ava.style.backgroundImage = `url(/towns/${town.type}.gif)`;
        if (!town.hasFort) {
            ava.style.backgroundImage = `url(/towns/${town.type}_0.png)`;
        }
        ava.style.borderColor = colors[report.color];
        lc.append(ava);

        let townInfo = document.createElement('table');
        townInfo.className = "townInfo";

        let row = document.createElement('tr');
        let c1 = document.createElement('td');
        let c2 = document.createElement('td');

        row.append(c1);
        row.append(c2);
        townInfo.append(row);

        let name = document.createElement('div');
        name.className = "townName value";
        name.innerText = town.name;
       
        let army = document.createElement('div');
        army.innerHTML = `Army: <span class='value'>${town.guardsValue}</span>&nbsp;&nbsp;`;

        lc.append(name);
        c1.append(army);

        c2.innerHTML = `<img class="regmIcon" src="/regm.png">`;
        let hasResearch = false;
        let lables = ["I", "II", "III", "IV", "V"];
        if (town.gmResearch) {
            for (let index = 0; index < 5; index++) {
                if (town.gmResearch[index] != 0) {
                    c2.innerHTML += `<span>${lables[index]}: <span class='value'>${town.gmResearch[index]}</span></span>&nbsp;`;
                    hasResearch = true;
                }
            }
        }

        if (!hasResearch) {
            c2.innerHTML = "";
        }

        row = document.createElement('tr');
        c1 = document.createElement('td');
        c2 = document.createElement('td');

        row.append(c1);
        row.append(c2);
        townInfo.append(row);


        for (let index = 0; index < town.spells.length; index++) {
            let spells = document.createElement('div');
            spells.className = `row spellsRow gm-${index+1}`;
            //spells.innerHTML = `${lables[index]}&nbsp;`;
            for (let n = 0; n <town.spells[index].length; n++) {
                const spell = town.spells[index][n];
                if (spell > 69 || spell <= 0) continue;
                spells.innerHTML += `<img class="spellIcon" src="/spells/${spell}.png">`;
            }
            if (index == 3 && town.spells.length == 5) {
                //spells.innerHTML = `&nbsp;&nbsp;${lables[index]}&nbsp;`;
                for (let n = 0; n <town.spells[4].length; n++) {
                    const spell = town.spells[4][n];
                    if (spell > 69 <= 0) continue;
                    spells.innerHTML += `<img class="spellIcon" src="/spells/${spell}.png">`;
                }
            }
            [c1, c2][index%2].append(spells);

            if (index%2 == 1) {
                row = document.createElement('tr');
                c1 = document.createElement('td');
                c2 = document.createElement('td');
        
                row.append(c1);
                row.append(c2);
                townInfo.append(row);
            }
        }

        
        item.append(townInfo);
        towns.append(item);
    });
}