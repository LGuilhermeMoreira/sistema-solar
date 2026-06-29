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
    // Shader de teste - torna tudo vermelho
    finalColor = vec4(1.0, 0.0, 0.0, 1.0);
}
