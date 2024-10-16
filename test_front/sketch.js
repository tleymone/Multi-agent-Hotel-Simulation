let mess = {};
let agt = {};
let hotel = {};
let day = 0;
let hour = 0;
let size = 20;
let roomsBox = 0;
let l = 0;

function setup() {
    createCanvas(windowWidth, windowHeight);
    //noStroke();
    frameRate(30);
}

function draw() {
    resizeCanvas(windowWidth + 440 * (roomsBox - 1), Math.max(windowHeight,560), true)
    
    background('#D9D9D9');
    textSize(size);
    fill('#ECE8E8');
    noStroke();
    rect(0, windowHeight*0.1, windowWidth + 440 * (roomsBox - 1), windowHeight * 0.9)
    fill('#000000');
    stroke(2);
    line(windowWidth * 0.3, windowHeight*0.1, windowWidth * 0.3, windowHeight);
    line(0, windowHeight*0.1, windowWidth + 440 * (roomsBox - 1), windowHeight*0.1);
    noStroke();
    text("Employés", 20, windowHeight*0.1 + size*1.4);
    text("Chambres", windowWidth * 0.3 + 20, windowHeight*0.1 + size*1.4);
    let api_url = 'http://localhost:8080/';
    roomsBox = 0;
    if (frameCount % 1 == 0) {
        httpDo(api_url = "data", "GET", "json", false, function (response) {
            day = response["time"].day;
            hour = response["time"].hour;
            mess = response["rooms"];
            agt = response["agents"];
            hotel = response["hotel"];
        });
    }
    if (mess[1] != undefined) {

    l = 0
    
        for (let i = 0; i < Object.keys(mess[1]).length; i++) {
            
            if (windowHeight*0.1 + 75 +size + size * 1.1  * i - (size * 1.1  * l) > windowHeight * 0.95) {
                roomsBox += 1;
                l = i;
            }
            text("Chambre " + mess[1][i].Number + ", capacity " + mess[1][i].Capacity + ", price " + mess[1][i].Price + ":",windowWidth * 0.3 + 20 + (440 * roomsBox), windowHeight*0.1 + 75 +size + size * 1.1  * i - (size * 1.1  * l));
            if (mess[1][i].State == 0) {
                fill('green');
            } else if (mess[1][i].State == 2){
                fill('orange')
            } else {
                fill('red');
            }
            ellipse(windowWidth*0.3 + 20*size + (440 * roomsBox), windowHeight*0.1 + 75 + (size * 0.6) + size * 1.1 * i - (size * 1.1  * l), size*0.75);
            fill(0);
        }
    }
    if (agt != undefined) {
        let nb = 0
        for (const i in agt) {
           //console.log(agt[i])
            text("Agt " + agt[i].agent.id + ":", 100, windowHeight*0.1 + 75 +size + size * 1.1 * nb);
            if (agt[i]["is-working"] == true) {
                fill('#22BB22');
            } else if (agt[i]["is-working"] == false){
                fill('#656565')
            }
            ellipse(13*size, windowHeight*0.1 + 75 + (size * 0.6) + size * 1.1 * nb, size*0.75);
            if (agt[i]["is-working"] == true) {
                if (agt[i].state == 0) {
                    fill(0);
                    text("libre", 16*size + 20, windowHeight*0.1 + 75 +size + size * 1.1 * nb);
                    fill('green');
                } else if (agt[i].state == 1){
                    fill(0);
                    text("travaille", 16*size + 20, windowHeight*0.1 + 75 +size + size * 1.1 * nb);
                    fill('red')
                }
                ellipse(16*size, windowHeight*0.1 + 75 + (size * 0.6) + size * 1.1 * nb, size*0.75);
            }
            if (agt[i].job == 0) {
                fill(0);
                text("rcpt", 35, windowHeight*0.1 + 75 +size + size * 1.1 * nb);
                fill('orange');
            } else if (agt[i].job == 1){
                fill(0);
                text("cln", 35, windowHeight*0.1 + 75 +size + size * 1.1 * nb);
                fill('blue')
            }
            ellipse(20, windowHeight*0.1 + 75 + (size * 0.6) + size * 1.1 * nb, size*0.75);

            fill(0);
            nb++
        }
    }
    if (hotel != undefined) {
        let str_money = "Argent : " + hotel.Money + " €";
       text(str_money, windowWidth/2 - size/2 * str_money.length/2, windowHeight*0.05 + size/2);   
    }
    text("Jour : " + day + ". Il est " + hour +"h.", 27, windowHeight*0.05 + size/2);
}