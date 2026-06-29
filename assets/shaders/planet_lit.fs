#version 330

in vec2 fragTexCoord;
in vec3 fragNormal;

uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform vec3 lightDir;

out vec4 finalColor;

void main()
{
    vec4 dayColor = texture(texture0, fragTexCoord)*colDiffuse;
    float sunAmount = max(dot(normalize(fragNormal), normalize(lightDir)), 0.0);
    float light = 0.12 + sunAmount*0.95;
    finalColor = vec4(dayColor.rgb*light, dayColor.a);
}
