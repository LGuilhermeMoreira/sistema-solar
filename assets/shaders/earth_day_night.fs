#version 330

in vec2 fragTexCoord;
in vec3 fragNormal;

uniform sampler2D texture0;
uniform sampler2D nightTexture;
uniform vec4 colDiffuse;
uniform vec3 lightDir;

out vec4 finalColor;

void main()
{
    // Normaliza a normal do fragmento
    vec3 norm = normalize(fragNormal);
    vec3 light = normalize(lightDir);
    
    // Calcula a quantidade de luz solar
    float sunAmount = dot(norm, light);
    
    // Suaviza a transição entre dia e noite
    float dayMix = smoothstep(-0.1, 0.2, sunAmount);
    
    // Amostra as texturas
    vec4 dayColor = texture(texture0, fragTexCoord) * colDiffuse;
    vec4 nightColor = texture(nightTexture, fragTexCoord);
    
    // Calcula a iluminação
    float daylight = max(sunAmount, 0.0) * 0.9 + 0.1;
    vec3 litDay = dayColor.rgb * daylight;
    vec3 nightGlow = nightColor.rgb * 0.8;
    
    // Mistura as cores de dia e noite
    vec3 finalRGB = mix(nightGlow, litDay, dayMix);
    
    finalColor = vec4(finalRGB, dayColor.a);
}
