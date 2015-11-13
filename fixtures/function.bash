v1 = "Global V1";
v2 = "Global V2";
function func(p) {
    v1 = "Local " ++ p;
    global v2;
    v2 = "V3 Modified.";
}
func("Var");
