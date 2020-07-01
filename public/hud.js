function hasItem(report, player, eid, flag) {
    player.querySelector(eid).className = report.hud[flag] ? "itemIcon hasItem" : "itemIcon";
}
function hasSpell(report, player, eid, flag) {
    player.querySelector(eid).className = report.hud[flag] ? "spellIcon hasItem" : "spellIcon";
}
let colors = ["red", "blue", "biege", "green", "orange", "purple", "teal", "pink"];

function fillHUD(player, report) {
    player.querySelector(".heroImage").style.backgroundImage = `url(/heroes/${report.startHero}.png)`;
    player.querySelector(".townImage").style.backgroundImage = `url(/towns/${report.startTown.type}.gif)`;    
    
    player.querySelector(".townImage").style.border = `4px solid ${colors[report.color]}`;
    player.querySelector(".heroImage").style.border = `4px solid ${colors[report.color]}`;

    player.querySelector("#playerName").innerText = report.playerName;

    player.querySelector("#tc .value").innerText = report.actions;
    player.querySelector("#apm .value").innerText = `${report.apm}/${report.cleanApm} (${report.maxApm})`;
    player.querySelector("#tmp .value").innerText = report.totalMPSpent;
    player.querySelector("#av .value").innerText = report.armyValue;

    hasItem(report, player, "#wings", "hasWings");
    hasItem(report, player, "#spellbinders_hat", "hasSpellbindersHat");
    hasItem(report, player, "#tome_of_air", "hasTomeOfAir");
    hasItem(report, player, "#tome_of_earth", "hasTomeOfEarth");
    hasItem(report, player, "#tome_of_fire", "hasTomeOfFire");

    hasItem(report, player, "#alliance", "hasAlliance");
    hasItem(report, player, "#aotd", "hasAOTD");
    hasItem(report, player, "#shackles", "hasShackles");

    hasSpell(report, player, "#Resurrect", "hasResurrect");
    hasSpell(report, player, "#dd", "hasDD");
    hasSpell(report, player, "#fly", "hasFly");
    hasSpell(report, player, "#tp", "hasTP");
    hasSpell(report, player, "#armageddon", "hasArmageddon");

    hasItem(report, player, "#grail", "hasGrail");
    hasItem(report, player, "#ee", "hasExpertEarth");

    player.querySelector(".combatInfo").style.display = report.currentCombat.heroArmyValue <= 0 ? "none" : "flex";
    player.querySelector(".smallHeroIcon").style.backgroundImage = `url(/heroes/${report.currentCombat.hero}.png)`;
    player.querySelector(".armyValues").innerHTML = `<span class="value">${report.currentCombat.heroArmyValue}</span> vs <span class="value">${report.currentCombat.enemyArmyValue}</span>`;

}